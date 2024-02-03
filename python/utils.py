from datetime import datetime
from pprint import pprint
import apsw
import dateutil.parser as dateparser
import feedparser
import os
import re
import requests
import time


def date2unix(timestamp):
    return int(time.mktime(dateparser.parse(timestamp, fuzzy=True).timetuple()))


def query_db(database, query, query_params):
    with apsw.Connection(database, flags=apsw.SQLITE_OPEN_READONLY) as db_connection:
        return db_connection.execute(query, query_params)


def cache_url(url, cache_path):
    resp = requests.get(
        url,
        headers={
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"
        },
    )

    if resp.status_code == 200:
        with open(cache_path, "w") as cache_file:
            cache_file.write(resp.text)

    elif resp.status_code == 403:
        pass

    else:
        raise Exception(resp.status_code)


def get_html(url, cache_dir, cache_time=300):
    if not os.path.exists(cache_dir):
        os.makedirs(cache_dir)

    cache_path = f"{cache_dir}/{url.replace('/', '_')}"

    if os.path.exists(cache_path):
        mtime = os.path.getmtime(cache_path)
        file_age = time.time() - mtime

        if file_age > cache_time:
            print("Requesting new file")
            cache_url(url, cache_path)

        else:
            print("Using cached file")

    else:
        print("Requesting new file")
        cache_url(url, cache_path)

    with open(cache_path, "r") as html_file:
        return html_file.read()


def get_feed_info(feedurl, cache_dir):
    html_file = get_html(feedurl, cache_dir)
    parsed_rss = feedparser.parse(html_file)["feed"]
    feed_info = {
        "rssurl": feedurl,
        "url": parsed_rss["link"],
        "title": parsed_rss["title"],
    }
    return feed_info


def get_entries(feedurl, cache_dir):
    html_file = get_html(feedurl, cache_dir)
    parsed_rss = feedparser.parse(html_file)
    return {i["id"]: i for i in parsed_rss["entries"]}


def parse_entry(feedurl, feed_entries, guid):
    if "author" in feed_entries[guid]:
        author = feed_entries[guid]["author"]
    else:
        author = ""

    return {
        "id": None,
        "guid": guid,
        "title": feed_entries[guid]["title"],
        "author": author,
        "url": feed_entries[guid]["link"],
        "feedurl": feedurl,
        "pubDate": date2unix(feed_entries[guid]["published"]),
        "content": feed_entries[guid]["summary"],
        "unread": 1,
        "enclosure_url": None,
        "enclosure_type": None,
        "deleted": 0,
    }


def test_query(*args, **kwargs):
    try:
        next(query_db(*args, **kwargs))
        return True

    except StopIteration:
        return False
