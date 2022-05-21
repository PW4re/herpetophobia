import  itertools
import collections
from attr import field
import websockets
import json
import asyncio
import random
from math import gcd
import  requests

URL = "0.0.0.0:5051"

S = SymmetricGroup(256)

#!/usr/bin/env python3

import enum
import random
import collections
from typing import List, Tuple, Set, Iterator, Sequence, Optional


Coordinates = Tuple[int, int]


class Direction(enum.Enum):
    LEFT = enum.auto()
    RIGHT = enum.auto()
    UP = enum.auto()
    DOWN = enum.auto()
    
    def opposite(self) -> 'Direction':
        if self is Direction.LEFT:
            return Direction.RIGHT

        if self is Direction.RIGHT:
            return Direction.LEFT

        if self is Direction.UP:
            return Direction.DOWN

        if self is Direction.DOWN:
            return Direction.UP
    
    def move(self, coordinates: Coordinates) -> Coordinates:
        x, y = coordinates

        if self is Direction.LEFT:
            return x - 1, y

        if self is Direction.RIGHT:
            return x + 1, y

        if self is Direction.UP:
            return x, y - 1

        if self is Direction.DOWN:
            return x, y + 1


class Snake:
    def __init__(self, initial: List[Coordinates], direction: Direction) -> None:
        self.deque = collections.deque(initial)
        self.direction = direction
        self.tail = self.deque[-1]
            
    @property
    def head(self) -> Coordinates:
        return self.deque[0]

    def copy(self) -> 'Snake':
        snake = Snake(self.deque, self.direction)
        snake.tail = self.tail

        return snake

    def has_intersection(self) -> bool:
        return self.deque.count(self.head) > 1

    def grow(self) -> None:
        self.deque.append(self.tail)
        
    def move(self, direction: Direction) -> None:
        if direction is not self.direction.opposite():
            self.direction = direction

        head = self.direction.move(self.head)
        self.deque.appendleft(head)

        self.tail = self.deque.pop()


MAX_STEPS = 256
LEVEL_WIDTH = 16
LEVEL_HEIGHT = 16
FOOD_COUNT = 8
FOOD_STEPS = LEVEL_WIDTH * LEVEL_HEIGHT // FOOD_COUNT


def get_possible_directions(snake: Snake) -> Iterator[Direction]:
    directions = [
        Direction.LEFT,
        Direction.RIGHT,
        Direction.UP,
        Direction.DOWN,
    ]

    for direction in directions:
        snake2 = snake.copy()

        if direction is snake.direction.opposite():
            continue

        snake2.move(direction)

        if snake2.has_intersection():
            continue

        x, y = snake2.head

        if not (0 <= x < LEVEL_WIDTH and 0 <= y < LEVEL_HEIGHT):
            continue

        yield direction


State = Tuple[Snake, List[Direction], int, int]


def find_path(
        state: State,
        target: Coordinates,
        visibility: Set[int],
        max_iterations: int,
        random_bound: int
) -> Optional[State]:
    queue: Sequence[State] = collections.deque([state])

    for _ in range(max_iterations):
        if len(queue) == 0:
            return

        current_snake, current_moves, current_index, steps = queue.pop()

        if steps > max(visibility):
            continue

        for direction in get_possible_directions(current_snake):
            if random.randint(0, random_bound) == 0:
                continue

            next_snake = current_snake.copy()
            next_snake.move(direction)

            next_moves = current_moves.copy()
            next_moves.append(direction)

            next_index = current_index

            if next_snake.head == target and steps in visibility:
                next_snake.grow()
                next_index += 1

                return next_snake, next_moves, next_index, steps + 1

            queue.append((next_snake, next_moves, next_index, steps + 1))


def solve(field: List[int], snake: Snake) -> List[Direction]:
    food: List[Coordinates] = [(0, 0)] * FOOD_COUNT

    for i, element in enumerate(field):
        x, y = i % LEVEL_WIDTH, i // LEVEL_HEIGHT

        if element % FOOD_STEPS == 0:
            position = element // FOOD_STEPS
            food[position] = (x, y)

    food_visibility: List[Set[int]] = []

    for i in range(len(food)):
        food_visibility.append(
            set([(i + 1) * FOOD_STEPS - 2, (i + 1) * FOOD_STEPS - 1]),
        )

    food = food
    food_visibility = food_visibility

    state: State = (snake, [], 0, 0)

    for target, visibility in zip(food, food_visibility):
        while True:
            next_state = find_path(state, target, visibility, 1_000, 2)

            if next_state is not None:
                state = next_state
                break

    return state[1]
    


def find_max_order(fields):
    return max(field.order() for field in fields)


def get_fields_by_counter(sign, log, offset, max_order, secret, init, counter):
    answer = None

    for t in range(-1000, 1000):
        if ((sign * (log + t) - offset * max_order) ^^ counter) % max_order == secret:
            answer = sign * (log + t)
            break

    return Permutation(init ^ (-1 * answer)), Permutation(init ^ answer)

def exploit(data, game_counter):
    counters = [counter for counter, _ in data]
    fields = [S(field) for _, field in data]
    
    max_order = find_max_order(fields)
    print(max_order)

    db = collections.defaultdict(list)
    bits = 7

    for f1, f2 in itertools.combinations(fields, 2):
        if f1.order() != max_order and f2.order() != max_order:
            continue

        try:
            log = discrete_log(f2, f1)
        except Exception:
            continue

        g, inv, _ = xgcd(log - 1, max_order)
        if g != 1:
            continue

        for t1 in range(1, 1_000):
            s = t1 * inv % max_order
            s_msb = (s >> bits) << bits

            db[s_msb].append(s)

        for t2 in range(1, 1_000):
            s = (-log * t2) * inv % max_order
            s_msb = (s >> bits) << bits

            db[s_msb].append(s)
    max_freq = max(len(ks) for ks in db.values())
    s_candidates = [s for s, ks in db.items() if len(ks) == max_freq]
    print(max_freq, s_candidates)

    init = None

    for f1, f2 in itertools.combinations(fields, 2):
        if f1.order() != max_order and f2.order() != max_order:
            continue

        try:
            log = discrete_log(f2, f1)
        except Exception:
            continue

        for s_msb in s_candidates:
            for t1, t2 in itertools.product(range(1, 2 ^ bits), repeat=2):
                if (s_msb + t1) * log % max_order == (s_msb + t2) % max_order:
                    t_diff = t1 - t2
                    g, t_diff_inv, _ = xgcd(t_diff, max_order)

                    if g > 1:
                        continue

                    init = (f1 / f2) ^ t_diff_inv
                    break

    print("INIT", Permutation(init))

    logs = []
    new_counters = []

    for COUNTER, FIELD in zip(counters, fields):
        try:
            log = discrete_log(FIELD, init)
        except Exception:
            continue

        logs.append(log)
        new_counters.append(COUNTER)

    print(len(logs), len(new_counters))

    secrets = set()

    for sign, offset in itertools.product([-1, 1], range(1, 2_000)):
        secrets.clear()

        for log, counter in zip(logs, new_counters):
            secrets.add(((sign * log - offset * max_order) ^^ counter) % max_order)

        if len(secrets) == 1:
            secret = secrets.pop()
            break

    print(offset, sign, secret)
    field1, field2 = get_fields_by_counter(sign, log, offset, max_order, secret, init, data[0][0])
    checkfield = Permutation(data[0][1])
    if checkfield == field1:
        valid_init = init
    elif checkfield == field2:
        valid_init = init ^ (-1)
    else:
        raise Exception("No valid init")

    next_field = get_fields_by_counter(sign, log, offset, max_order, secret, valid_init, game_counter)
    return next_field[0]

def multiply(left_perm, right_perm):
    new_perm = [0 for _ in range(len(left_perm))]
    for i, element in enumerate(right_perm):
        new_perm[i] = left_perm[element]
    
    return new_perm


def invert(perm):
    new_perm = [0 for _ in range(len(perm))]
    for i in range(len(perm)):
        new_perm[perm[i]] =  i
    
    return new_perm


def exponentiation(perm, power):
    result = [i for i in range(len(perm))]

    if power < 0:
        perm = invert(perm)
        power = -power

    while power > 0:
        if power & 1 == 1:
            result = multiply(result, perm)
        
        perm = multiply(perm, perm)
        power  = power >> 1
    
    return result

def create_game():
    init_perm = [157, 79, 170, 8, 108, 234, 163, 16, 251, 181, 23, 148, 55, 162, 211, 186, 194, 222, 152, 207, 57, 97, 87, 45, 245, 141, 142, 40, 13, 92, 89, 64, 191, 102, 247, 178, 28, 138, 118, 68, 226, 24, 151, 103, 15, 139, 154, 244, 180, 83, 82, 196, 171, 167, 31, 155, 63, 246, 38, 200, 228, 120, 218, 204, 10, 238, 47, 56, 146, 185, 172, 158, 133, 53, 117, 42, 193, 241, 206, 86, 161, 0, 77, 243, 149, 239, 121, 129, 2, 85, 159, 59, 96, 164, 81, 220, 114, 18, 214, 65, 60, 125, 188, 201, 104, 174, 153, 75, 240, 223, 126, 35, 189, 113, 27, 236, 122, 143, 124, 73, 227, 43, 49, 67, 187, 48, 99, 250, 39, 20, 165, 115, 1, 177, 93, 232, 202, 249, 116, 54, 6, 242, 252, 69, 255, 22, 176, 197, 110, 5, 61, 169, 254, 183, 19, 229, 109, 150, 111, 131, 156, 253, 208, 145, 58, 179, 76, 7, 91, 78, 37, 233, 212, 9, 215, 192, 62, 209, 33, 32, 198, 168, 17, 195, 136, 166, 98, 130, 71, 248, 90, 217, 25, 30, 112, 34, 231, 3, 237, 21, 80, 224, 100, 66, 52, 84, 106, 4, 101, 205, 26, 105, 128, 225, 210, 135, 137, 175, 95, 70, 132, 203, 182, 29, 219, 190, 199, 44, 235, 140, 147, 74, 144, 46, 123, 216, 221, 14, 94, 127, 119, 36, 184, 88, 107, 12, 41, 134, 213, 72, 173, 160, 50, 51, 11, 230]
    order = 8933296680
    power = random.randint(1, order - 1)
    while gcd(order, power) != 1:
        power = random.randint(1, order - 1)
    secret = str(random.randint(1000, 893329668))
    print(secret)
    body = {
        "secret": secret,
        "init": exponentiation(init_perm, power),
        "flag": "FLAGGG" + str(random.randint(1, 10000))
    }

    res = requests.post(f"http://{URL}/create", json=body)
    return res.json()["id"]

def print_map(map):
    for el in map:
        print(" ".join(el))

async def sploit(id):
    counters_fields = []
    async with websockets.connect(f"ws://{URL}/play") as ws:
        for i in range(30):
            await ws.send(str(json.dumps({"id": id})))
            resp = await ws.recv()
            resp = json.loads(resp)
            while not resp.get("gameResult"):
                await ws.send(str(json.dumps({"direction": "w",
                    "closeGame": False,
                    "newGame": False})))
                resp = await ws.recv()
                resp = json.loads(resp)
            counters_fields.append((resp["counter"], [x + 1 for x in resp["permutation"]]))
            if i != 29:
                await ws.send(str(json.dumps({
                    "direction": "w",
                    "closeGame": False,
                    "newGame": True
                })))
                resp = await ws.recv()
                resp = json.loads(resp)
        await ws.send(str(json.dumps({
            "direction": "",
            "closeGame": False,
            "newGame": True
        })))
        resp = await ws.recv()
        resp = json.loads(resp)
        res = exploit(counters_fields, resp["counter"])
        res = [x - 1 for x in res]
        print("field for counter", resp["counter"], "is", res)
        snake = Snake([(4, 1), (3, 1), (2, 1), (1, 1)], Direction.RIGHT)
        directions = solve(res, snake)

        for direction in directions:
            if direction == Direction.UP:
                move = "w"
            elif direction == Direction.DOWN:
                move = "s"
            elif direction == Direction.LEFT:
                move = "a"
            elif direction == Direction.RIGHT:
                move = "d"
            else:
                raise Exception("Bad direction!!")
            await ws.send(str(json.dumps({
                "direction": move,
                "closeGame": False,
                "newGame": False
            })))
            resp = await ws.recv()
            print(direction)
            try:
                print_map(json.loads(resp)["gameMap"])
                print("=============================")
            except Exception:
                print(json.loads(resp))
        
        # while not resp.get("gameResult"):
        #         await ws.send(str(json.dumps({"direction": "w",
        #             "closeGame": False,
        #             "newGame": False})))
        #         resp = await ws.recv()
        #         resp = json.loads(resp)
        # print(resp["permutation"])

# websocket.enableTrace(True)
asyncio.get_event_loop().run_until_complete(sploit(create_game()))
