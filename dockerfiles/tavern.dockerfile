FROM python:3.12-slim
ARG TAVERN_FILE="./acceptance-tests/tavern/test_config.tavern.yaml"

RUN groupadd -r test \
    && useradd -rg test test -m

USER test:test

WORKDIR /home/test

COPY ${TAVERN_FILE} ./test.tavern.yaml

COPY requirements.txt ./requirements.txt
RUN pip install -r ./requirements.txt

ENTRYPOINT ["python", "-m", "pytest"]
CMD ["./test.tavern.yaml"]
