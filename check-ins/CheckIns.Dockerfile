FROM python:3.10-slim

WORKDIR /opt/app

COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY . /opt/app/
ENV PYTHONPATH=/opt/app/

ENTRYPOINT ["uvicorn", "controller:app", "--host", "0.0.0.0", "--port", "8080"]
