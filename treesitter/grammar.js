module.exports = grammar({
	name: 'nvimboat',

	rules: {
		nvimboat: $ => repeat($._element),
		_element: $ => choice(
			// $._tags_page,
			$._filter,
			$._feed,
			$._article,
			$.header,
			$.body,
		),
		_filter: $ => choice(
			$.unread_filter,
			$.read_filter
		),
		_feed: $ => choice(
			$.unread_feed,
			$.read_feed
		),
		_article: $ => choice(
			$.unread_article,
			$.read_article
		),
		unread_filter: $ => / \| N \(\d+\/\d+\).*? \| query:.*?, tags:.*?/,
		read_filter: $ => / \|   \(\d+\/\d+\).*? \| query:.*?, tags:.*?/,

		unread_feed: $ => / \| N \(\d+\/\d+\).*? \| http.*?/,
		read_feed: $ => / \|   \(\d+\/\d+\).*? \| http.*?/,

		unread_article: $ => / \| N \| \d\d (?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec).*?/,
		read_article: $ => / \|   \| \d\d (?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec).*?/,

		header: $ => seq($._header_line),
		_header_line: $ => seq($._header_key, ': ', $._expression),
		_header_key: $ => choice('Feed', 'Title', 'Author', 'Date', 'Link'),

		body: $ => seq(
			$.body_start,
			$.content,
		),
		body_start: $ => /== Article Begin ==/,
		content: $ => /(.*?\n)+?/,

		_expression: $ => /.*/,
	}
});
