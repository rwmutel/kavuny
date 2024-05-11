FROM python:3.10-slim

# RUN apt install libpq-dev

WORKDIR /opt/app

COPY requirements.txt requirements.txt
RUN pip install --upgrade pip
RUN pip install -r requirements.txt

COPY ./coffee_pack_service.py /opt/app/
COPY ./controller.py /opt/app/
COPY ./persistence.py /opt/app/
COPY ./coffee_pack_model.py /opt/app/

ENTRYPOINT ["uvicorn", "controller:app", "--host", "0.0.0.0", "--port", "8080"]
