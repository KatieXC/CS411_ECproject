# CS411: Extra Credit Project
## Running the Program
~~~
go run joins.go <input filename 1> <join column name in file 1> <input filename 2> <join column name in file 2> <join method (HASH or NESTED_LOOP)> <output filename>
~~~

For example, running:
~~~ 
go run joins.go nation.csv REGIONKEY region.csv REGIONKEY HASH out.csv
~~~
will perform a join between nation and region tables on the attribute regionkey and write its output to out.csv.  The time taken to perform the join will also be printed to the screen.

Note: The column names are case sensitive, but the join method is not.
## Input Files
Due to the size of some tables, the join operation could take an extra long time. To decrease the table size without changing the input file data, change the variable STEPLENGTH on line 14 in the *joins.go* file to only access every *x* rows.

Ex: STEPLENGTH = 10 --> Only read in one row for every 10 rows.

The input files' delimeters can also be updated in the *joins.go* file on line 36.

## Benchmarks
These benchmarks were collected using input tables lineitem and orders and join column ORDERKEY.  In order to get enough benchmarks in a reasonable time, I set STEPLENGTH to 20 (see above section for how).

| Hash Join         || Nested Loop Join     ||
| Run | Time (ms)    | Run | Time ()         |
|:---:|:------------:|:---:|:---------------:|
|  1  | 130.688833ms |  1  | 3m20.260940667s |
|  2  | 132.420791ms |  2  | 3m16.470997959s |
|  3  | 104.526959ms |  3  | 3m24.418210541s |
|  4  | 111.628292ms |  4  | 3m16.802281666s |
|  5  | 127.158583ms |  5  | 3m19.824742s    |
|  6  | 112.298666ms |  6  | 3m23.553532375s |
|  7  | 113.63225ms  |  7  | 3m18.3828045s   |
|  8  | 121.411042ms |  8  | 3m13.970427542s |
|  9  | 124.892333ms |  9  | 3m17.088949417s |
|  10 | 107.913708ms |  10 | 3m21.263024875s |
| **AVG** | 118.6571457ms| **AVG** | 3m19.203591154s |

Based on these results, the hash join algorithm is over 1600 times faster on average than the nested loop algorithm!  There is a *significant* performance improvement when the nested for loop is avoided.