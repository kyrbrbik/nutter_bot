from debian:11-slim
COPY requirements.txt /tmp/requirements.txt
RUN apt-get update && apt-get --no-install-recommends install -y python3 python3-pip 
RUN pip3 install -r /tmp/requirements.txt
COPY main.py /app/main.py
WORKDIR /app
CMD ["python3", "main.py"]
