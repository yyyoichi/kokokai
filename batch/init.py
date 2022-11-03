from src.db import get_connection


def init():
    conn = get_connection()
    cursor = conn.corsor()
    with open("init.sql", "r", encoding="utf-8") as f:
        sql = f.read()
        cursor.execute(sql)
        conn.commit()
    conn.close()
