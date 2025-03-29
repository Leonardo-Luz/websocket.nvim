local M = {}

local job
local ensure_job = function()
  if job then
    return job
  end

  job = vim.fn.jobstart({ 'go/websocket.nvim' }, { rpc = true })
  return job
end

local run_goolang_func = function(func_name, ...)
  vim.fn.rpcrequest(ensure_job(), func_name, ...)
end

M.start_ws_client = function()
  run_goolang_func('newWsClient', "localhost", "8080")
  -- run_goolang_func('connectToWsServer')
end

vim.api.nvim_create_user_command('StartWsClient', function()
  M.start_ws_client()
end, {})

return M
