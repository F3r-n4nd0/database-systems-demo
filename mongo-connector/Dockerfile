FROM python:3.7

RUN pip install mongo-connector mongo-connector[solr] solr_doc_manager

RUN mkdir /conf && mkdir /var/log/mongo-connector
COPY ./wait-for-it.sh .
RUN chmod +x /wait-for-it.sh