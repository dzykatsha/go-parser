#!/usr/bin/env bash
not_exist_file_cmd=(go run client.go -i not_exist.html -o output.txt)
wrong_extension_cmd=(go run client.go -i i_am.wrong -o also.wrong)
forgotten_arg_cmd=(go run client.go)
for FILE in $1/* 
do 
    FULL_FILENAME=$FILE
    FILENAME=${FULL_FILENAME##*/}
    OUTPUT=${FILENAME%%.*}
    html_cmd=(go run client.go -i $FILE -o $2/${OUTPUT}_html.txt)
    doc_cmd=(go run client.go -i $FILE -o $2/${OUTPUT}_doc.txt)
    docx_cmd=(go run client.go -i $FILE -o $2/${OUTPUT}_docx.txt)
    pdf_plain_cmd=(go run client.go -i $FILE -o $2/${OUTPUT}_pdf_plain.txt --plain)
    pdf_grouped_cmd=(go run client.go -i $FILE -o $2/${OUTPUT}_pdf_grouped.txt --grouped)
    pdf_styled_cmd=(go run client.go -i $FILE -o $2/${OUTPUT}_pdf_styled.txt --styled)
    wrong_flag_cmd=(go run client.go -i $FILE -o $2/${OUTPUT}_pdf_not_exist.txt --notexist)

    case "$FILE" in
        *.html) "${html_cmd[@]}" & wait  ;;
        *.pdf)  "${pdf_plain_cmd[@]}" & wait; "${pdf_grouped_cmd[@]}" & wait; "${pdf_styled_cmd[@]}" & wait; "${wrong_flag_cmd[@]}" & wait ;;
        *.doc)  "${doc_cmd[@]}" & wait  ;;
        *.docx) "${docx_cmd[@]}" & wait ;;
    esac
done
"${not_exist_file_cmd[@]}" & wait
"${wrong_extension_cmd[@]}" & wait
"${forgotten_arg_cmd[@]}" & wait
