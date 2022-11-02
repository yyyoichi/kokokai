import psycopg2
import src.setting as setting


class DB:
    def __init__(self) -> None:
        self._connect()
        print("Connect db")

    def _connect(self):
        self.db = psycopg2.connect(
            user=setting.DB_USER,
            password=setting.DB_PASS,
            host=setting.DB_HOST,
            port=setting.DB_PORT,
            database=setting.DB_NAME,
            sslmode="verify-ca",
            sslrootcert=setting.SSL_ROOT_CERT,
            sslcert=setting.SSL_CERT,
            sslkey=setting.SSL_KEY
        )

    def test(self):
        cur = self.db.cursor()
        cur.execute("select version()")
        print(cur.fetchone())

        cur.close()

    def close(self):
        self.db.close()
        print("Close db\n")
