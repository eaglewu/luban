#! /bin/bash

if [ $# -lt 1 ]; then
    echo "usage: ./compare.sh directory"
    exit 1
fi

if [ ! -d "$1" ]; then
    echo "no such directory"
    exit 1
fi

find "$1" -type f -name "*.php" | while read fname; do
    # Genrate native token
    if [ ! -e "$fname.json" ]; then
        php -r "
            if(! extension_loaded('Tokenizer')){
                echo \"Tokenizer extension is required.\", PHP_EOL;
                exit(1);
            }
            \$scripts=file_get_contents('$fname');
            \$results = [];    
            \$tokens = token_get_all(\$scripts);
            foreach(\$tokens as \$token){
                if (is_array(\$token)){
                    \$tokenName = token_name(\$token[0]);
                    if (\$tokenName == 'T_WHITESPACE') {
                        continue;
                    }
                    \$results[] = [
                        'type'  =>  1,
                        'l'     =>  \$token[2],
                        't'     =>  \$tokenName,
                        'v'     =>  \$token[1],
                    ];
                } else {
                    \$results[] = [
                        'type'  =>  2,
                        'v'     =>  \$token,
                    ];
                }
            }
            echo json_encode(\$results, JSON_UNESCAPED_UNICODE);
        " > "$fname.json"
    fi

    if [ ! $? -eq 0 ]; then
        echo "Failed to genrate native tokens[$fname]"
        exit 1
    fi

    $GOPATH/bin/lexer -file="$fname" -compare="$fname.json"

    if [ ! $? -eq 0 ]; then
        echo -e "[\033[31mFAILED\033[0m] $fname"
        exit 1
    fi
done

if [ $? -eq 0 ]; then
    echo 
    echo "Congratulations!"
    echo "[$1] pass the tests."
else
    exit 1
fi
