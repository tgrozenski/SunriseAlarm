import json
import time
import boto3
import jwt
import requests
from botocore.exceptions import ClientError

SNS_TOPIC_ARN = "arn:aws:sns:us-west-1:136595888751:observability-topic"
sns_client = boto3.client("sns", region_name="us-west-1")

def get_secret(secret_name: str, region: str = "us-west-1") -> dict:
    client = boto3.client("secretsmanager", region_name=region)
    try:
        response = client.get_secret_value(SecretId=secret_name)
        return json.loads(response["SecretString"])
    except ClientError as e:
        raise e

def get_access_token(service_account: dict) -> str:
    now = int(time.time())

    payload = {
        "iss": service_account["client_email"],
        "scope": "https://www.googleapis.com/auth/firebase.messaging",
        "aud": service_account["token_uri"],
        "iat": now,
        "exp": now + 3600
    }

    signed_jwt = jwt.encode(
        payload,
        service_account["private_key"],
        algorithm="RS256"
    )

    response = requests.post(
        service_account["token_uri"],
        data={
            "grant_type": "urn:ietf:params:oauth:grant-type:jwt-bearer",
            "assertion": signed_jwt
        }
    )
    response.raise_for_status()
    return response.json()["access_token"]

def send_fcm_message(access_token: str, project_id: str, topic: str, title: str, body: str):
    url = f"https://fcm.googleapis.com/v1/projects/{project_id}/messages:send"

    message = {
        "message": {
            "topic": topic,
            "notification": {
                "title": title,
                "body": body
            }
        }
    }

    response = requests.post(
        url,
        headers={
            "Authorization": f"Bearer {access_token}",
            "Content-Type": "application/json"
        },
        json=message
    )
    response.raise_for_status()
    return response.json()

def notify(subject: str, message: dict):
    sns_client.publish(
        TopicArn=SNS_TOPIC_ARN,
        Subject=subject,
        Message=json.dumps(message, indent=2)
    )

def handler(event, context):
    try:
        secret = get_secret("LAMBDA_SERVICE_KEY")
        access_token = get_access_token(secret)

        result = send_fcm_message(
            access_token=access_token,
            project_id=secret["project_id"],
            topic="your-topic-name",
            title="Hello",
            body="This is a test notification"
        )

        response = {
            "statusCode": 200,
            "body": json.dumps(result)
        }
        notify("FCM Success", response)
        return response

    except Exception as e:
        response = {
            "statusCode": 500,
            "body": json.dumps({
                "error": str(e),
                "type": type(e).__name__
            })
        }
        notify("FCM Failure", response)
        return response

