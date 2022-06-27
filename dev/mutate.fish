#! /usr/bin/env fish

# Mutate. Based loosely on:
# * https://mutmut.readthedocs.io/en/latest/
# * https://github.com/zimmski/go-mutesting

# Would eventually like to replace with a go program with its own testing
# Would like to be able to say "replace bool returns with their opposites"
# Would like to cache candidates and results

# true -> false
set search_text true
set replacement false

set files (ag $search_text -G '.*/.go$' -l)
set command 'go mod tidy &&
  golangci-lint run -c ./dev/golangci.toml --fix 2> /dev/null &&
  go test -rapid.nofailfile -failfast &&
  ./fuzz.fish'

set caught true

for file in $files
    echo mutatable '"'$search_text'"' found in $file
    set candidates (ag --column $search_text $file)
    for candidate in $candidates
        set line (echo $candidate | awk -F: '{print $1}')
        set column (math (echo $candidate | awk -F: '{print $2}')-1)
        set match (echo $candidate | awk -F: '{print $3}')
        set mutant (echo $match | sed -E 's/(.{'$column'})'$search_text'/\1'$replacement'/')
        echo mutating $file:$line:(math $column+1) '"'$match'"' '->' '"'$mutant'"'
        sed -i "" -E $line's/(.{'$column'})'$search_text'/\1'$replacement'/' $file
        if eval $command > /dev/null
            echo failed to catch the mutant
            set caught false
        else
            echo caught the mutant!
        end
        echo restoring mutant $file:$line:(math $column+1) '"'$mutant'"' '->' '"'$match'"'
        sed -i "" -E $line's/(.{'$column'})'$replacement'/\1'$search_text'/' $file
        if test $caught = false
            return 1
        end
    end
end

if test -z $files
    echo no files were found with candidates for mutation
end

return 0
