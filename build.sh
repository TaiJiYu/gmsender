#!/bin/bash
OutDir=tool
RUN_NAME=gmsender.exe
rm -f ${OutDir}/${RUN_NAME}

rsrc -ico pkg/asset/favicon.ico -o rsrc.syso
go build -ldflags="-s -w -H windowsgui" -trimpath -o ${OutDir}/${RUN_NAME}
rm -f rsrc.syso  # 清理临时文件

${OutDir}/${RUN_NAME}