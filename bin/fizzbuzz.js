#!/usr/bin/env node

/**
 * Enhanced FizzBuzz implementation
 * Prints numbers from 1 to 100 in ascending order, but:
 * - For multiples of 3, print 'Fizz' instead of the number
 * - For multiples of 5, print 'Buzz' instead of the number
 * - For multiples of both 3 and 5, print 'FizzBuzz'
 * - Additionally, prints the prime factors of each number at the end of the line
 */

/**
 * Returns the prime factors of a number as an array
 * @param {number} num - The number to factorize
 * @returns {number[]} - Array of prime factors
 */
function getPrimeFactors(num) {
    if (num <= 1) return [num];
    
    const factors = [];
    let divisor = 2;
    
    while (num > 1) {
        if (num % divisor === 0) {
            factors.push(divisor);
            num /= divisor;
        } else {
            divisor++;
        }
    }
    
    return factors;
}

function customFizzBuzz(start = 1, end = 100) {
    for (let i = start; i <= end; i++) {
        let output = '';
        
        if (i % 3 === 0 && i % 5 === 0) {
            output = 'FizzBuzz';
        } else if (i % 3 === 0) {
            output = 'Fizz';
        } else if (i % 5 === 0) {
            output = 'Buzz';
        } else {
            output = i.toString();
        }
        
        // Get and display prime factors
        const primeFactors = getPrimeFactors(i);
        const factorsStr = primeFactors.join(' × ');
        
        console.log(`${output} [${factorsStr}]`);
    }
}

// Execute the custom FizzBuzz when this script is run directly
if (require.main === module) {
    console.log("Running FizzBuzz from 1 to 100:");
    customFizzBuzz(1, 100);
}

// Export the function for potential reuse
module.exports = customFizzBuzz;