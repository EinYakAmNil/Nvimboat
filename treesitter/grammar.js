module.exports = grammar({
	name: 'nvimboat',

	rules: {
		nvimboat: $ => repeat($._element),
		_element: $ => choice(
			$._feed,
			$._article,
			$.header,
			$.body,
		),

		_feed: $ => choice(
			$.unread_feed,
			$.read_feed
		),
		_article: $ => choice(
			$.unread_article,
			$.read_article
		),

		unread_feed: $ => / \| N \(\d+\/\d+\).*?/,
		read_feed: $ => / \|   \(\d+\/\d+\).*?/,

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
