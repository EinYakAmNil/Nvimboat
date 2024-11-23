local M = {}

function M:push(page_type, page_id)
	self[#self + 1] = { type = page_type, id = page_id }
end

function M:pop()
	self[#self] = nil
end

return M
