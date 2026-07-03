"""Fetch Blue Archive JP Excel resources.

绕过 asfu222/BlueArchiveLocalizationTools 的 il2cpp dump 步骤(新版游戏的
libil2cpp.so 被改名并加固,Il2CppInspector 无法处理)。改用 Deathemonic/BA-FB
仓库 data 分支定期 CI 产出的 types.cs(dump.cs)直接生成 FlatData schema,再用
BLT 的 TableExtractor 解析 Excel.zip;ExcelDB.db 现在用 SQLCipher 加密,用
sqlcipher3 + JP 服的 SQLCipher key 解密后,同样用 BLT 处理。

产出 ./resources/Excel/*.json 与 ./resources/ExcelDB/*.json,供 generate_excel.go
打包成 data/Excel.bin。
"""
import json
import os
import re
import shutil
import sys
from pathlib import Path

BLT_DIR = Path("./BlueArchiveLocalizationTools")
sys.path.insert(0, str(BLT_DIR))

from extractor import compile_python  # noqa: E402

RESOURCES = Path("./resources")
TABLE_BUNDLES = RESOURCES / "TableBundles"
EXTRACTED = Path("./Extracted")
TYPES_CS = EXTRACTED / "Dumps" / "types.cs"
FLATDATA_DIR = EXTRACTED / "FlatData"
TABLE_OUT = EXTRACTED / "Table"

JP_SQLCIPHER_KEY = os.environ.get(
    "BAPS_JP_SQLCIPHER_KEY",
    "ef0aaca06f34b4a4be3172a75a3ea565e815f9ece35b1fb12b7a166ba0807bc4",
)


def patch_blt_compiler_bug(flatdata_dir: Path) -> None:
    """修复 BLT compiler.py 生成的 GachaSelectPickupGroupExcel.py 重复定义。

    BLT 的 compiler 把同名 struct 的重复定义拼接进同一文件,导致 builder 方法
    带错误缩进(IndentationError)。这些方法只用于构建 flatbuffer,不影响读取,
    直接删掉重复段即可。
    """
    target = flatdata_dir / "GachaSelectPickupGroupExcel.py"
    if not target.exists():
        return
    text = target.read_text(encoding="utf8")
    fixed = re.sub(
        r"\n def AddNameKr\(builder.*\Z",
        "\n",
        text,
        flags=re.DOTALL,
    )
    if fixed != text:
        target.write_text(fixed, encoding="utf8")
        print(f"patched {target.name}")


def generate_flatdata() -> None:
    """从 types.cs 生成 FlatData python schema(含 dump_wrapper.py)。"""
    if FLATDATA_DIR.exists():
        print("FlatData already exists, skipping generation")
        return
    print("Compiling FlatData from types.cs...")
    compile_python(str(TYPES_CS), str(EXTRACTED))
    patch_blt_compiler_bug(FLATDATA_DIR)
    print(f"FlatData generated: {len(os.listdir(FLATDATA_DIR))} files")


def extract_excel_zip() -> None:
    """用 BLT TableExtractor 提取 Excel.zip(双层加密:zip 密码 + 字段 XOR)。"""
    print("Extracting Excel.zip...")
    from xtractor.table import TableExtractor

    extractor = TableExtractor(
        str(TABLE_BUNDLES), str(TABLE_OUT), "Extracted.FlatData"
    )
    extractor.extract_zip_file("Excel.zip")
    print("Excel.zip done")


def extract_exceldb_db() -> None:
    """用 sqlcipher3 解密 ExcelDB.db 并用 BLT 处理每行 .bytes。

    ExcelDB.db 用 SQLCipher 整库加密;解密后每张表(如 CharacterDBSchema)的 Bytes
    列存的是单条 <Name>Excel flatbuffer(BLT 的 _process_bytes_file 会按非 Table
    路径走 dump_<Name>)。每张表聚合后写出 <Name>Excel.json。
    """
    print("Extracting ExcelDB.db...")
    import sqlcipher3

    from xtractor.table import TableExtractor

    extractor = TableExtractor(
        str(TABLE_BUNDLES), str(TABLE_OUT), "Extracted.FlatData"
    )

    db_path = TABLE_BUNDLES / "ExcelDB.db"
    if not db_path.exists():
        raise FileNotFoundError(f"missing {db_path}")

    dst = TABLE_OUT / "ExcelDB"
    dst.mkdir(parents=True, exist_ok=True)

    conn = sqlcipher3.connect(str(db_path))
    try:
        conn.execute(f"PRAGMA key = \"x'{JP_SQLCIPHER_KEY}'\"")
        conn.execute("SELECT count(*) FROM sqlite_master").fetchone()
        table_names = [
            r[0]
            for r in conn.execute(
                "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"
            )
        ]
        print(f"found {len(table_names)} tables in ExcelDB.db")

        ok = fail = 0
        for schema_name in sorted(table_names):
            excel_name = schema_name.replace("DBSchema", "Excel")
            rows = conn.execute(f'SELECT Bytes FROM "{schema_name}"').fetchall()
            if not rows:
                continue
            merged = []
            for (blob,) in rows:
                if not blob:
                    continue
                result, _ = extractor._process_bytes_file(excel_name, blob)
                if result:
                    if isinstance(result, list):
                        merged.extend(result)
                    else:
                        merged.append(result)
            if merged:
                (dst / f"{excel_name}.json").write_text(
                    json.dumps(merged, indent=4, ensure_ascii=False), encoding="utf8"
                )
                ok += 1
            else:
                fail += 1
        print(f"ExcelDB done: ok={ok}, fail={fail}")
    finally:
        conn.close()


def move_to_resources() -> None:
    """把 Extracted/Table/{Excel,ExcelDB} 移到 resources/,供 generate_excel.go 读取。"""
    for name in ("Excel", "ExcelDB"):
        src = TABLE_OUT / name
        dst = RESOURCES / name
        if not src.exists():
            raise FileNotFoundError(f"missing {src}")
        if dst.exists():
            shutil.rmtree(dst)
        shutil.move(str(src), str(dst))
        print(f"moved {src} -> {dst}")


def main() -> None:
    generate_flatdata()
    extract_excel_zip()
    extract_exceldb_db()
    move_to_resources()
    print("All resources fetched.")


if __name__ == "__main__":
    main()
