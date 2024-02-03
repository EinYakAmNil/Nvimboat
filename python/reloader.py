import apsw
from utils import date2unix, get_html, query_db, test_query, get_feed_info, get_entries, parse_entry


class Updater:
    guid_query = "SELECT guid FROM rss_item WHERE feedurl = :feedurl"
    feedurl_query = "SELECT rssurl FROM rss_feed WHERE rssurl = :feedurl"
    insert_rss_item = "INSERT into rss_item values(:id, :guid, :title, :author, :url, :feedurl, :pubDate, :content, :unread, :enclosure_url, :enclosure_type, 0, '', :deleted, '', 'text/html', '', '')"
    insert_rss_feed = "INSERT into rss_feed values(:rssurl, :url, :title, 0, 0, '')"

    def __init__(self, database: str, cache_dir: str, *rss_url: str):
        self.database = database
        self.cache_dir = cache_dir
        self.rss_urls = rss_url
        self.new_feeds = self._check_rssurls()
        self.entries = {url: get_entries(url, self.cache_dir) for url in self.rss_urls}
        self.new_rss_items = sum(
            [self._get_new_items(url) for url in self.rss_urls], start=()
        )

    def _get_new_items(self, feedurl):
        feed_entries = self.entries[feedurl]

        known_guids = [
            i[0]
            for i in query_db(
                self.database, self.guid_query, {"feedurl": feedurl}
            )
        ]
        new_guids = (set(known_guids) ^ set(feed_entries)) - set(known_guids)

        return tuple(parse_entry(feedurl, feed_entries, guid) for guid in new_guids)

    def _check_rssurls(self) -> list:
        return tuple(
            get_feed_info(url, self.cache_dir)
            for url in self.rss_urls
            if not test_query(
                self.database, self.feedurl_query, {"feedurl": url}
            )
        )

    def update_database(self):
        with apsw.Connection(
            self.database, flags=apsw.SQLITE_OPEN_READWRITE
        ) as db_connection:
            db_connection.executemany(self.insert_rss_feed, self.new_feeds)
            db_connection.executemany(self.insert_rss_item, self.new_rss_items)


if __name__ == "__main__":
    feedurls = [
        "https://suckless.org/atom.xml",
        "https://www.archlinux.org/feeds/news/",
        "https://lukesmith.xyz/rss.xml",
        "https://www.pathofexile.com/news/rss",
    ]
    feedurl = "https://lukesmith.xyz/rss.xml"

    database = "cache/cache.db"
    cache_dir = "cache/"
    updater = Updater(database, cache_dir, *feedurls)
    # updater = Updater(database, cache_dir, feedurl)
    updater.update_database()
