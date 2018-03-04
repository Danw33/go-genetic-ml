/**
 * go-genetic-ml
 *
 * A Golang-based Genetic Machine Learning Algorithm
 *
 * Written by Daniel Wilson (@Danw33) <hello@danw.io>
 * With special thanks to the book "The Nature of Code" by Daniel Shiffman
 *
 * https://github.com/Danw33/go-genetic-ml
 *
 * @copyright Copyright (C) 2018 Daniel J. Wilson <hello@danw.io>
 * @license GNU GPL v3.0 - See LICENSE
 */
/**
  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

/**
 * Adjustable Variables
 */
var (
	// Target Outcome
	target = "I think, therefore I am."

	// Maximum Popultaion
	maxpop = 250

	// Mutation Rate
	mutrate float32 = 0.01
)

/**
 * DNA
 * Represents a single entity, there genes (rune slice) and assessed fitness
 */
type DNA struct {
	genes   []rune
	fitness float32
}

/**
 * Population
 * Holds the entities of the population, the mating pool, and iteration information
 */
type Population struct {
	entities     []DNA
	matingPool   []DNA
	generations  int
	completed    bool
	perfectScore float32
}

/**
 * Main Method
 * Sets up the initial generation, then runs the evolution loop until an entity
 * matches the target 100%
 */
func main() {
	fmt.Println("Danw33's Golang-based Genetic Algorithm")
	fmt.Println("Start time:", time.Now())
	fmt.Println("Running with Max Population:", maxpop, "and Mutation Probability:", mutrate)
	fmt.Println("Target Outcome: ", target)

	// Seed the PRNG (only once!)
	rand.Seed(time.Now().Unix())

	// Sanity Check
	//test()

	var population = Population{[]DNA{}, []DNA{}, 0, false, 1.0}

	// Run the setup method (Create Generation 0)
	setup(&population)

	// Evolve
	for population.completed == false {
		evolve(&population)
	}

	fmt.Println("Solution Discovered at", time.Now(), "by Generation", population.generations, "with population", len(population.entities), "and mutation rate", mutrate, " Average fitness:", populationAverageFitness(&population), "Final Phrase:", populationGetBest(&population))
}

/**
 * Initial Setup Method
 * Generates Generation 0 of the population with all-new DNA (Random)
 */
func setup(population *Population) {
	fmt.Println("Setting up at", time.Now())

	fmt.Println("Populating Generation 0 Gene Pool with random DNA Geonomes")
	for i := 0; i < maxpop; i++ {
		var newDna = DNA{}
		dnaCreate(&newDna, len(target))
		population.entities = append(population.entities, newDna)
	}

	fmt.Println("Created Seed Entities:", len(population.entities))

	fmt.Println("Calculating Generation 0 Fitness")
	populationCalculateFitness(population, target)
	fmt.Println("Generation 0 Fitness has been calculated.")

	fmt.Println("Setup Completed at", time.Now())
}

/**
 * Evolution Loop Method
 * Runs the Natural Selection, Generation, Fitness cycle
 * To be called in a loop until the population flags itself as completed.
 */
func evolve(population *Population) {
	// Generate mating pool
	populationNaturalSelection(population)

	// Create next generation
	populationGenerate(population)

	// Calculate fitness
	populationCalculateFitness(population, target)

	// Display Info
	fmt.Println("Generation", population.generations, "with population", maxpop, "and mutation rate", mutrate, "completed with average fitness", populationAverageFitness(population), "Best Phrase:", populationGetBest(population))

}

func test() {

	fmt.Println("Running basic test. Will Generate two parents, crossover and mutuate.")

	var dnaA = DNA{}
	dnaCreate(&dnaA, len(target))
	dnaAssessFitness(&dnaA, target)
	fmt.Println("Parent 1 (DNA A) Fitness:", dnaA.fitness, "Phrase:", dnaExtractPhrase(&dnaA))

	var dnaB = DNA{}
	dnaCreate(&dnaB, len(target))
	dnaAssessFitness(&dnaB, target)
	fmt.Println("Parent 2 (DNA B) Fitness:", dnaB.fitness, "Phrase:", dnaExtractPhrase(&dnaB))

	var dnaC = dnaCrossover(&dnaA, &dnaB)
	dnaMutate(&dnaC, mutrate)
	dnaAssessFitness(&dnaC, target)
	fmt.Println("Child    (DNA C) Fitness:", dnaC.fitness, "Phrase:", dnaExtractPhrase(&dnaC))

	fmt.Println("Manipulating Child geonome (DNA C => DNA D) to test fitness assessment")

	var dnaD = DNA{}
	var mutatedGenes []rune
	mutatedGenes = append(mutatedGenes, rune(target[0])) // Mutate the gene at the position 0
	mutatedGenes = append(mutatedGenes, rune(target[1])) // Mutate the gene at the position 1
	mutatedGenes = append(mutatedGenes, rune(target[2])) // Mutate the gene at the position 2
	mutatedGenes = append(mutatedGenes, dnaC.genes[3:]...)
	dnaD.genes = mutatedGenes

	dnaAssessFitness(&dnaD, target)
	fmt.Println("Child    (DNA D) Fitness:", dnaD.fitness*100, "Phrase:", dnaExtractPhrase(&dnaD))

	fmt.Println("Testing concluded, see console for data to analyse.")
}

/**
 * Random Int Generator with Range Restriction
 * Generates a random int within the given min and max parameters
 * Uses math/rand library
 */
func random(min, max int) int {
	return rand.Intn(max-min) + min
}

/**
 * Random Float Generator with Range Restriction
 * Generates a random float within the given min and max parameters
 * Uses math/rand library
 */
func randomFloat(min, max float32) float32 {
	return rand.Float32()*(max-min) + min
}

/**
 * Re-maps a number from one range to another
 * Danw33's Golang interpretation of the Arduino Math Library map method
 * @see https://www.arduino.cc/reference/en/language/functions/math/map/
 */
func highLowMap(x, inMin, inMax, outMin, outMax float32) float32 {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

/**
 * DNA: Create New, Random DNA
 * Creates n new DNA genes,
 * Appends them to the genes array (rune slice) in the given dna struct pointer
 */
func dnaCreate(dna *DNA, n int) {
	for i := 0; i < n; i++ {
		dna.genes = append(dna.genes, rune(random(32, 128))) // Pick from range of chars
	}
}

/**
 * DNA: Extract the genes as a string
 * Built from the genes rune slice in the given dna pointer
 */
func dnaExtractPhrase(dna *DNA) string {
	return string(dna.genes)
}

/**
 * DNA: Fitness Assessment Method
 * Sets a percentage (float32) of "correct" runes (how close to the target) on
 * the given dna pointer
 */
func dnaAssessFitness(dna *DNA, target string) {
	var score int
	var runeTarget = []rune(target)

	for i := 0; i < len(dna.genes); i++ {
		if dna.genes[i] == runeTarget[i] {
			score++
		}
	}

	dna.fitness = float32(score) / float32(len(target))
}

/**
 * DNA: Crossover Method
 * Takes two DNA Parents, and returns a DNA Child that has genes spliced from
 * both parents
 */
func dnaCrossover(partnerA *DNA, partnerB *DNA) DNA {
	// Create a new child
	var child = DNA{}

	// Pick a midpoint in the genes
	var midpoint = random(0, len(partnerA.genes))

	// Half from one, half from the other
	for i := 0; i < len(partnerA.genes); i++ {
		if i > midpoint {
			// Before the midpoint, take partner A's genes
			// In Java: child.genes[i] = partnerA.genes[i];
			child.genes = append(child.genes, partnerA.genes[i])
		} else {
			// After the midpoint, take partner B's genes
			child.genes = append(child.genes, partnerB.genes[i])
		}
	}

	// Return the new child
	return child
}

/**
 * DNA: Mutation Method
 * Mutates the genes of the given entity, within the given mutation rate (probability)
 */
func dnaMutate(entity *DNA, rate float32) {
	for i := 0; i < len(entity.genes); i++ {
		if randomFloat(0.0, 1.0) < rate {
			// In Java: genes[i] = (char) random(32,128);
			var mutatedGenes []rune
			mutatedGenes = append(mutatedGenes, entity.genes[:i]...)   // NB: append() is a variadic function
			mutatedGenes = append(mutatedGenes, rune(random(32, 128))) // Mutate the gene at the random position
			mutatedGenes = append(mutatedGenes, entity.genes[i+1:]...) // the ... lets you us multiple arguments to a variadic function from a slice
			entity.genes = mutatedGenes
		}
	}
}

/**
 * Population: Run a fitness assessment on every current member of the population
 */
func populationCalculateFitness(population *Population, target string) {
	for i := 0; i < len(population.entities); i++ {
		dnaAssessFitness(&population.entities[i], target)
	}
}

/**
 * Population: Mating Pool Generator
 * Performs Natural Selection on the current generation of entities, and creates
 * a mating pool of DNA candidates to become parents.
 */
func populationNaturalSelection(population *Population) {
	// Reset the mating pool first
	population.matingPool = []DNA{}

	var maxFitness float32

	// Find the fittest entity in the current population
	for i := 0; i < len(population.entities); i++ {
		if population.entities[i].fitness > maxFitness {
			maxFitness = population.entities[i].fitness
		}
	}

	// Each member of the current population will be added to the new mating pool a given number of times
	// based on their assessed fitnes. The higher the fitness, the more entries a single entity will have
	// therefore increasing the chances of a fitter child being produced (Natural Selection)
	for i := 0; i < len(population.entities); i++ {
		var fitness = highLowMap(population.entities[i].fitness, 0, maxFitness, 0, 1)
		var n = int(fitness * 100) // Like the book we use an Arbitrary multiplier. An alternative would be the monte carlo method.
		for j := 0; j < n; j++ {   // We pick out two random numbers
			population.matingPool = append(population.matingPool, population.entities[i])
		}
	}
}

/**
 * Population: Generation Iteration
 * Replaces the population's entities with the new entities generated
 * from the mating pool, performing DNA crossover and mutation.
 */
func populationGenerate(population *Population) {
	// Refill the population with children from the mating pool
	for i := 0; i < len(population.entities); i++ {
		var a, b int
		a = int(random(0, len(population.matingPool)))
		b = int(random(0, len(population.matingPool)))

		var partnerA, partnerB, child DNA
		partnerA = population.matingPool[a]
		partnerB = population.matingPool[b]
		child = dnaCrossover(&partnerA, &partnerB)

		dnaMutate(&child, mutrate)
		population.entities[i] = child
	}

	population.generations++
}

/**
 * Population: Get Best
 * Gets the best phrase generated by the entity of the current population with
 * the highest fitness (here known as the "world record")
 */
func populationGetBest(population *Population) string {
	var worldrecord float32
	var index int

	for i := 0; i < len(population.entities); i++ {
		if population.entities[i].fitness > worldrecord {
			index = i
			worldrecord = population.entities[i].fitness
		}
	}

	if worldrecord == population.perfectScore {
		population.completed = true
	}

	return dnaExtractPhrase(&population.entities[index])
}

/**
 * Population: Average Fitness
 * Calculates and returns the average fitness for the current generation of
 * the population
 */
func populationAverageFitness(population *Population) float32 {
	var total float32
	for i := 0; i < len(population.entities); i++ {
		total += population.entities[i].fitness
	}
	return total / float32(len(population.entities))
}

/**
 * Population: All Phrases
 * Currently unused method, capable of outputting all phrases held by each entity
 * within the current population. Can be called within the evolution loop to help
 * with debugging.
 */
func populationAllPhrases(population *Population) string {
	var everything string
	var displayLimit int = int(math.Min(float64(len(population.entities)), 50))

	for i := 0; i < displayLimit; i++ {
		everything += dnaExtractPhrase(&population.entities[i]) + "\n"
	}

	return everything
}
