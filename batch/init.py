from src.db import get_connection


def init():
    conn = get_connection()
    cursor = conn.cursor()
    with open("init.sql", "r", encoding="utf-8") as f:
        sql = f.read()
        cursor.execute(sql)
        conn.commit()
    conn.close()


if __name__ == "__main__":
    init()
