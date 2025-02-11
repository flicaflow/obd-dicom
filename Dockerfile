FROM alpine:latest

RUN mkdir /app

WORKDIR /app

ADD build/linux/amd64/obd-dicom /app

ENTRYPOINT [ "/app/obd-dicom", "-scp", "-calledae", "DICOM_SCP", "-port", "1040", "-datastore", "/datastore" ]
