local M = {}

M.path = {}

function M.reset()
	M.path = {}
	return M.path
end

function M.append(p)
	M.path[#M.path+1] = p
	return M.path
end

function M.subtract()
	M.path[#M.path] = nil
	return M.path
end

return M
