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

vim.api.nvim_create_user_command('StartWsClient', function(args)
  local set_var = vim.api.nvim_buf_set_var
  local get_var = vim.api.nvim_buf_get_var
  local buf = vim.api.nvim_get_current_buf()
  set_var(buf, "is_ws_update", false)

  M.start_ws_client(args.args)

  vim.api.nvim_create_autocmd("TextChangedI", {
    callback = function()
      if get_var(buf, "is_ws_update") then
        set_var(buf, "is_ws_update", false)
        return
      end

      local lines = vim.api.nvim_buf_get_lines(buf, 0, -1, false);
      local line = table.concat(lines, "Ef232wefeEFAwdEFF")

      run_goolang_func('writeWsClient', 0, -1, line, current_url)
    end
  })
end, { nargs = 1 })

return M
