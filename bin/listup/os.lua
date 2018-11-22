-- os, distribution, version, bit

-- OSの取得
local os_type = exec("uname")
setValue("$.os.os", os_type)

if os_type == "Linux" then 
  -- ディストリビューションとバージョンの取得
  -- ディストリビューションは/etc/xxxx-releaseファイルがあるかで判断

  dist = "unknown"

  -- Ubuntu, Debian
  if exist("/etc/debian_version") or exist("/etc/debian_release") then
    if exist("/etc/lsb-release") then
      dist = "Ubuntu"
    else
      dist = "Debian"
    end

  -- Fedora
  elseif exist("/etc/fedora-release") then
    dist = "Fedora"

  -- RedHat, Oracle, CentOS
  elseif exist("/etc/redhat-release") then
    if exist("/etc/centos-release") then
      dist = "CentOS"
    elseif exist("/etc/oracle-release") then
      dist = "Oralce"
    else
      dist = "RedHat Enterprise"
    end

  -- Arch
  elseif exist("/etc/arch-release") then
    dist = "Arch Linux"

  -- Turbo
  elseif exist("/etc/turbo-release") then
    dist = "Turbo Linux"

  -- SuSE
  elseif exist("/etc/SuSE-release") then
    dist = "Open SuSE"

  -- Vine
  elseif exist("/etc/vine-release") then
    dist = "Vine Linux"

  -- Gentoo
  elseif exist("/etc/gentoo-release") then
    dist = "Gentoo Linux"
  end

  setValue("$.os.ditribution", dist)

  setValue("$.os.version", "")

  -- Bit数(32/64)
  setValue("$.os.bit", exec("getconf LONG_BIT"))
end

