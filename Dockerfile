FROM python

RUN pip install discord.py

copy in discord.py

copy in k8s.py

get env vars for discord and k8s connection

run discord.py, passing in connection variables.

ENV ADMIN_PASS="temp"

ENTRYPOINT [ "sh", "-c", "python discord.py $ADMIN_PASS"]
