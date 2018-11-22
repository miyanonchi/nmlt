function exec(prog)
  local handle = io.popen(prog)
  local result = handle:read()
  handle:close()

  return result
end

function exist(path)
  if type(path) ~= "string" then
    return false
  end

  if is_dir(path) then
    return true
  end

  local ok, err, code = os.rename(path, path)

  if not ok then
    if code == 13 then
      -- Permission denied, but it exists
      return true
    end

    return false
  end

  return true
end

function is_dir(path)
  if type(path) ~= "string" then
    return false
  end

  local res = os.execute("cd " .. path .. " 2> /dev/null")

  if res == nil then
    return false
  end

  return true
end
