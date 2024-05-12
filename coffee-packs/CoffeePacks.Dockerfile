FROM python:3.10-slim

# RUN apt install libpq-dev

WORKDIR /opt/app

COPY requirements.txt requirements.txt
RUN pip install --upgrade pip
RUN pip install -r requirements.txt

COPY *.py /opt/app/

ENTRYPOINT ["uvicorn", "controller:app", "--host", "0.0.0.0", "--port", "8080"]
