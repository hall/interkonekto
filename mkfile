install:
 rsync --delete -lr . linguo.io:/srv/ido/api --exclude .git

backup:V:
 sqlite3 ido.db '.dump' > ido.sql

restore:V:
 rm ido.db; cat ido.sql | sqlite3 ido.db
