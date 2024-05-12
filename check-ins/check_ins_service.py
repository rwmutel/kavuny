from persistence import Client

conn = Client()


def get_check_ins(coffee_shop_id: int | None,
                  user_id: int | None,
                  coffee_pack_id: int | None):
    return conn.get_check_ins(coffee_shop_id, user_id, coffee_pack_id)


def post_check_ins(check_in: dict):
    return conn.post_check_in(check_in)
