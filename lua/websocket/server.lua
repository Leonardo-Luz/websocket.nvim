local M = {}

local job
local ensure_job = function()
  if job then
    return job
  end

  local script_path = vim.fn.stdpath 'data' .. "../../go/websocket.nvim"

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
  print("Running start_ws_server...") -- Log when this function is invoked
  run_goolang_func('startWsServer', vim.api.nvim_buf_get_lines(vim.api.nvim_get_current_buf(), 0, -1, false), port)
end

vim.api.nvim_create_user_command('StartWsServer', function(args)
  M.start_ws_server(args.args)
end, { nargs = 1 })

return M
