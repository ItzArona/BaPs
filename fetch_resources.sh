#!/usr/bin/env bash
# Fetch Blue Archive JP Excel resources.
#
# 新版游戏(1.70+)的 libil2cpp.so 被改名并加固,asfu222/BlueArchiveLocalizationTools
# 的 setup_flatdata.py 会因找不到 libil2cpp.so 而失败,且 Il2CppInspector 无法
# dump 加固后的 il2cpp。这里改用 Deathemonic/BA-FB 仓库 data 分支定期 CI 产出的
# types.cs 直接生成 FlatData schema,跳过 il2cpp dump。ExcelDB.db 现在用
# SQLCipher 加密,用 sqlcipher3 + JP 服 key 解密。
set -euo pipefail

# 1. BlueArchiveLocalizationTools(BLT):提供 compile_python / TableExtractor
if [ ! -d BlueArchiveLocalizationTools ]; then
  git clone --depth=1 https://github.com/asfu222/BlueArchiveLocalizationTools
fi
pip3 install -q -r BlueArchiveLocalizationTools/requirements.txt
pip3 install -q sqlcipher3

# 2. types.cs:BA-FB data 分支定期 dump 的 il2cpp C# 类型定义(绕过加固 il2cpp)
mkdir -p Extracted/Dumps
if [ ! -f Extracted/Dumps/types.cs ]; then
  echo "Downloading types.cs from Deathemonic/BA-FB data branch..."
  curl -fL --retry 3 -o Extracted/Dumps/types.cs \
    "https://raw.githubusercontent.com/Deathemonic/BA-FB/data/Japan/cs/types.cs"
fi

# 3. 拿当前 catalog URL,下载 Excel.zip / ExcelDB.db
python3 BlueArchiveLocalizationTools/update_urls.py ba.env ./data/ServerInfo.json
export $(grep -v '^#' ba.env | xargs)
echo "Using catalog url: $ADDRESSABLE_CATALOG_URL"
mkdir -p resources/TableBundles
curl -fL --retry 3 "${ADDRESSABLE_CATALOG_URL}/TableBundles/Excel.zip" \
  -o resources/TableBundles/Excel.zip
curl -fL --retry 3 "${ADDRESSABLE_CATALOG_URL}/TableBundles/ExcelDB.db" \
  -o resources/TableBundles/ExcelDB.db

# 4. 生成 FlatData schema + 提取 Excel.zip + 解密并提取 ExcelDB.db
python3 fetch_resources.py
