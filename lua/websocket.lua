require 'websocket.server'
require 'websocket.client'

local M = {}

M.setup = function() end

return M

-- -- IGNORE
-- M.print_go = function()
--   local buf = vim.api.nvim_get_current_buf()
--   local win = vim.api.nvim_get_current_win()
--   local current_line = vim.api.nvim_win_get_cursor(win)[1] - 1
--
--   local lines = { "hello", "from", "Goolang" }
--
--   run_goolang_func('hello', buf, current_line, current_line + #lines, lines)
-- end
--
-- vim.api.nvim_create_user_command('Hello', function()
--   M.print_go()
-- end, {})
