FROM public.ecr.aws/lambda/python:3.14

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY src/main.py ${LAMBDA_TASK_ROOT}/

CMD ["main.handler"]
