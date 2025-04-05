local M = {}

local job
local ensure_job = function()
  if job then
    return job
  end

  local script_path = vim.fn.stdpath 'data' .. "/lazy/websocket.nvim/go/websocket.nvim"

  job = vim.fn.jobstart({ script_path }, { rpc = true })
  return job
end

local run_goolang_func = function(func_name, ...)
  vim.fn.rpcrequest(ensure_job(), func_name, ...)
end

local current_url = nil

M.start_ws_client = function(url)
  local buf = vim.api.nvim_get_current_buf();
  current_url = url
  run_goolang_func('startWsClient', buf, current_url)
end

-- local processed = false

vim.api.nvim_create_user_command('StartWsClient', function(args)
  M.start_ws_client(args.args)

  vim.api.nvim_create_autocmd("InsertLeave", {
    callback = function()
      -- if processed then
      --   processed = false
      --   return
      -- end

      local line = vim.api.nvim_get_current_line();
      -- if line:find("^EWFOIJO324wefEFWWEFFjnksdv09") then
      --   local parsed = line:gsub("^EWFOIJO324wefEFWWEFFjnksdv09", "")
      --   vim.print(line)
      --   vim.api.nvim_set_current_line(parsed)
      --   processed = true
      --   return
      -- end

      -- processed = true


      local cursor = vim.api.nvim_win_get_cursor(vim.api.nvim_get_current_win());
      run_goolang_func('writeWsClient', cursor[1], cursor[1], line, current_url)
    end
  })
end, { nargs = 1 })

return M
