local action = require("nvimboat.action")
local M = {}
local map_opts = { silent = true, buffer = 0, nowait = true }

M.keymaps = {
	n = {
		q = {
			rhs = action.show_main_menu,
			opts = map_opts
		},
		l = {
			rhs = action.select,
			opts = map_opts
		},
		h = {
			rhs = action.back,
			opts = map_opts
		},
		t = {
			rhs = action.show_tags,
			opts = map_opts
		},
		a = {
			rhs = action.toggle_article_read,
			opts = map_opts
		},
		n = {
			rhs = action.next_unread,
			opts = map_opts
		},
		p = {
			rhs = action.prev_unread,
			opts = map_opts
		},
		N = {
			rhs = action.prev_unread,
			opts = map_opts
		},
		J = {
			rhs = action.next_article,
			opts = map_opts
		},
		K = {
			rhs = action.prev_article,
			opts = map_opts
		},
		o = {
			rhs = action.open_media,
			opts = map_opts
		},
		R = {
			rhs = action.reload_all,
			opts = map_opts
		},
		D = {
			rhs = action.delete,
			opts = map_opts
		},
	},
	v = {
		a = {
			rhs = action.toggle_article_read,
			opts = map_opts
		},
		o = {
			rhs = action.open_media,
			opts = map_opts
		},
		n = {
			rhs = action.next_unread,
			opts = map_opts
		},
		p = {
			rhs = action.prev_unread,
			opts = map_opts
		},
		N = {
			rhs = action.prev_unread,
			opts = map_opts
		},
		D = {
			rhs = action.delete,
			opts = map_opts
		},
	},
}

function M.configure(opts)
	if opts.keymaps then
		for mode, keymap in pairs(opts.keymaps) do
			for lhs, map in pairs(keymap) do
				M.keymaps[mode][lhs] = { rhs = map.rhs, opts = map.opts }
			end
		end
	end

	return M.keymaps
end

return M
