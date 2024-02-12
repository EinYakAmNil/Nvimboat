#!/bin/env python3
import apsw
import argparse
import feedparser
import logging
import os
import requests
import sys
import time
import dateutil.parser as dateparser

FEEDURL_QUERY = "SELECT rssurl FROM rss_feed"
GUID_QUERY = "SELECT guid FROM rss_item"
INSERT_RSS_FEED = (
    "INSERT INTO rss_feed(rssurl, url, title, etag) VALUES(:rssurl, :url, :title, '')"
)
INSERT_RSS_ITEM = "INSERT INTO rss_item(guid, title, author, url, feedurl, pubDate, content, unread, content_mime_type) VALUES(:guid, :title, :author, :url, :feedurl, :pubDate, :content, 1, 'text/html')"


def date2unix(timestamp):
    return int(time.mktime(dateparser.parse(timestamp, fuzzy=True).timetuple()))


def cache_url(url: str, cache_path: str) -> bool:
    resp = requests.get(
        url,
        headers={
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36"
        },
    )
    with open(cache_path, "w") as cache_file:
        cache_file.write(resp.text)
    if resp.status_code == 200:
        return True
    else:
        logging.warning(f"{url} had error: {resp.status_code}")
        return False


def get_html(url: str, cache_dir: str, cache_time: int = 300) -> str:
    if not os.path.exists(cache_dir):
        os.makedirs(cache_dir)

    cache_path = f"{cache_dir}/{url.replace('/', '_')}"

    if os.path.exists(cache_path):
        mtime = os.path.getmtime(cache_path)
        file_age = time.time() - mtime

        if file_age > cache_time:
            logging.info(f"Requesting new file for: {url}")
            if not cache_url(url, cache_path):
                logging.warning(f"Error caching file for: {url}")

        else:
            logging.info(f"Using cached file for: {url}")

    else:
        logging.info(f"Requesting new file for: {url}")
        if not cache_url(url, cache_path):
            logging.warning(f"Error caching file for: {url}")

    try:
        with open(cache_path, "r") as html_file:
            return html_file.read()
    except:
        return ""


def parse_feed(feedurl: str, cache_dir: str, cache_time: int):
    try:
        html = get_html(feedurl, cache_dir, cache_time)

    except Exception as e:
        logging.warning(f"Error getting HTML for: {feedurl}.\n{e}\n{type(e)}")
        return None, None

    try:
        parsed_rss = feedparser.parse(html)
        feed_info = {
            "rssurl": feedurl,
            "url": parsed_rss["feed"]["link"],
            "title": parsed_rss["feed"]["title"],
        }
        entries = {i["id"]: i for i in parsed_rss["entries"]}
        return feed_info, entries

    except Exception as e:
        logging.warning(f"Error parsing content for: {feedurl}.\n{e}\n{type(e)}")
        return None, None


def parse_entry(feedurl, guid, entry):
    if "author" in entry:
        author = entry["author"]
    else:
        author = ""

    try:
        return {
            "guid": guid,
            "title": entry["title"],
            "author": author,
            "url": entry["link"],
            "feedurl": feedurl,
            "pubDate": date2unix(entry["published"]),
            "content": entry["summary"],
        }
    except Exception as e:
        logging.warning(f"Error parsing entry: {guid}.\n{e}")
        return None


def feed_generator(feedurl: str, cache_dir: str, cache_time: int):
    feed_info, entries = parse_feed(feedurl, cache_dir, cache_time)
    if feed_info and entries:
        entries = [
            parse_entry(feedurl, guid, entry)
            for guid, entry in entries.items()
            if parse_entry(feedurl, guid, entry) is not None
        ]
        return feed_info, entries
    else:
        return None


def update_db(database: str, feeds: str, articles: str):
    db_conn = apsw.Connection(database, flags=apsw.SQLITE_OPEN_READWRITE)
    known_feeds = {rssurl[0] for rssurl in db_conn.execute(FEEDURL_QUERY)}
    known_articles = {guid[0] for guid in db_conn.execute(GUID_QUERY)}
    new_feeds = (known_feeds ^ set(feeds)) - known_feeds
    new_articles = (known_articles ^ set(articles)) - known_articles
    for f in new_feeds:
        logging.info(f"New feed: {f}")
    for a in new_articles:
        logging.info(f"New article: {a}")
    db_conn.executemany(
        INSERT_RSS_FEED, [f for k, f in feeds.items() if k in new_feeds]
    )
    logging.info("Inserted new feeds.")
    db_conn.executemany(
        INSERT_RSS_ITEM, [a for k, a in articles.items() if k in new_articles]
    )
    logging.info("Inserted new articles.")
    db_conn.close()
    logging.info("Closed database.")


if __name__ == "__main__":
    cmdline = argparse.ArgumentParser(description="Reloader for Nvimboat feeds")
    cmdline.add_argument(
        "-d",
        "--cache-dir",
        default="./cache",
        type=str,
        help="Directory for the database and the cached requests.",
    )
    cmdline.add_argument(
        "-t",
        "--cache-time",
        default=300,
        type=int,
        help="Duration in seconds for which cached files are considered valid.",
    )
    cmdline.add_argument(
        "-v",
        "--verbose",
        action="store_true",
        help="",
    )
    cmdline.add_argument("urls", nargs="+", help="URLs to request.")
    args = vars(cmdline.parse_args())
    cache_dir = args["cache_dir"]
    database = cache_dir + "/cache.db"
    cache_time = args["cache_time"]
    feedurls = args["urls"]
    if args["verbose"]:
        loglevel = logging.INFO
    else:
        loglevel = logging.FATAL

    logging.basicConfig(
        level=loglevel,
        format="%(levelname)s: %(message)s",
        handlers=[logging.FileHandler(f"{cache_dir}/pyboat.log"), logging.StreamHandler(sys.stdout)],
    )

    parsed_feeds = [feed_generator(f, cache_dir, cache_time) for f in feedurls]
    feeds = {
        entry[0]["rssurl"]: entry[0] for entry in parsed_feeds if entry is not None
    }
    articles = {
        entry["guid"]: entry
        for entry_list in parsed_feeds
        if entry_list is not None
        for entry in entry_list[1]
    }
    update_db(database, feeds, articles)
