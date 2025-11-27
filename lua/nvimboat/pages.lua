---@alias PageType
---| "Article"
---| "Feed"
---| "Filter"
---| "MainMenu"
---| "TagFeeds"
---| "TagsOverview"

---@class Page
---@field type PageType
---@field id string

---@class PageStack
---@field [integer] Page
local M = {}

---@param page_type PageType
---@return nil
function M:push(page_type, page_id)
	self[#self + 1] = { type = page_type, id = page_id }
end

---@return Page
function M:pop()
	local tail = self[#self]
	self[#self] = nil
	return tail
end

---@return nil
function M:reset()
	for p in ipairs(self) do
		self[p] = nil
	end
end

return M
