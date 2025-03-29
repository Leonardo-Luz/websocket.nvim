local M = {}

local job
local ensure_job = function()
  if job then
    return job
  end

  -- Start the Go WebSocket server process
  print("Starting Go WebSocket server...") -- Log when the job is being started
  job = vim.fn.jobstart({ 'go/websocket.nvim' }, { rpc = true })
  return job
end

local run_goolang_func = function(func_name, ...)
  local job_id = ensure_job()
  if job_id then
    vim.fn.rpcrequest(job_id, func_name, ...)
  end
end

M.start_ws_server = function()
  print("Running start_ws_server...") -- Log when this function is invoked
  run_goolang_func('newWsServer', "localhost", "8080")
end

vim.api.nvim_create_user_command('StartWsServer', function()
  M.start_ws_server()
end, {})

return M
