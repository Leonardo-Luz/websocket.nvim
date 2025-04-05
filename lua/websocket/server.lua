local M = {}

local job
local ensure_job = function()
  if job then
    return job
  end

  local script_path = vim.fn.stdpath 'data' .. "/lazy/websocket.nvim/go/websocket.nvim"

  -- Start the Go WebSocket server process
  print("Starting Go WebSocket server...") -- Log when the job is being started
  job = vim.fn.jobstart({ script_path }, { rpc = true })
  return job
end

local run_goolang_func = function(func_name, ...)
  local job_id = ensure_job()
  if job_id then
    vim.fn.rpcrequest(job_id, func_name, ...)
  end
end

M.start_ws_server = function(port)
  local bufnum = vim.api.nvim_get_current_buf()

  print("Running start_ws_server...") -- Log when this function is invoked
  run_goolang_func('startWsServer', vim.api.nvim_buf_get_lines(bufnum, 0, -1, false), port)
  vim.defer_fn(function()
    run_goolang_func('startWsClient', bufnum, "ws://localhost:" .. port .. "/ws")

    local set_var = vim.api.nvim_buf_set_var
    local get_var = vim.api.nvim_buf_get_var
    local buf = vim.api.nvim_get_current_buf()
    set_var(buf, "is_ws_update", false)

    vim.api.nvim_create_autocmd("TextChangedI", {
      callback = function()
        if get_var(buf, "is_ws_update") then
          set_var(buf, "is_ws_update", false)
          return
        end

        local lines = vim.api.nvim_buf_get_lines(buf, 0, -1, false);
        local line = table.concat(lines, "Ef232wefeEFAwdEFF")

        run_goolang_func('writeWsClient', 0, -1, line, "ws://localhost:" .. port .. "/ws")
      end
    })
  end, 100)
end

vim.api.nvim_create_user_command('StartWsServer', function(args)
  M.start_ws_server(args.args)
end, { nargs = 1 })

return M
