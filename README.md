# GoTSP

* **Input:** CSV file of n-dimensional points
* **Output:** total Euclidean distance with optimal path (permutation)
* **Methodology:** Genetic algorithm
  * **Reproduction strategy:** Partially mapped crossover (PMX)
  * **Mutation:** Subpath reversal

GoTSP solves the traveling salesperson problem using a genetic algorithm. A set of points can be imported from a CSV file. The shortest permutation will be written out to a new or existing CSV file which can then be used again as input. Each time the program runs, a random population of a given size is generated with the imported shortest path. The mutation function is the reversal if a subpath in a member of the population.