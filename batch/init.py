import sys
from src.db import get_connection


def delete():
    print("delete rows")
    conn = get_connection()
    cursor = conn.cursor()
    with open("sql/delete_row.sql", "r", encoding="utf-8") as f:
        sql = f.read()
        cursor.execute(sql)
        conn.commit()
    conn.close()


def init():
    conn = get_connection()
    cursor = conn.cursor()
    with open("sql/init.sql", "r", encoding="utf-8") as f:
        sql = f.read()
        cursor.execute(sql)
        conn.commit()
    conn.close()


def user_init():
    conn = get_connection()
    cursor = conn.cursor()
    with open("sql/user.sql", "r", encoding="utf-8") as f:
        sql = f.read()
        cursor.execute(sql)
        conn.commit()
    conn.close()


def test():
    conn = get_connection()
    cursor = conn.cursor()
    cursor.execute('SELECT CURRENT_TIMESTAMP')
    print(cursor.fetchone())
    conn.close()


if __name__ == "__main__":
    arg = sys.argv[1]
    print(arg)
    if arg == "init":
        init()
    elif arg == "delete":
        delete()
    elif arg == "test":
        test()
    elif arg == "user":
        user_init()
