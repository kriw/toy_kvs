for i in $(seq 5)
do
    python -c "print('A' * 2 * 1024 * 1024)" > ./files/twoMega$i.txt
    echo $i >> ./files/twoMega$i.txt
done

