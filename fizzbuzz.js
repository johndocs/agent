/**
 * Enhanced FizzBuzz implementation with prime factorization
 * - Prints numbers from 1 to 100 in ascending order
 * - For multiples of 3, print "Fizz" instead of the number
 * - For multiples of 5, print "Buzz" instead of the number
 * - For multiples of both 3 and 5, print "FizzBuzz"
 * - Each line ends with the prime factorization of the number
 */

/**
 * Finds the prime factors of a number
 * @param {number} num - The number to factorize
 * @returns {string} String representation of prime factors
 */
function getPrimeFactors(num) {
  if (num <= 1) return num.toString();
  
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
  
  // Group repeated factors with exponents
  const factorMap = {};
  factors.forEach(factor => {
    factorMap[factor] = (factorMap[factor] || 0) + 1;
  });
  
  // Format the factors as a string
  return Object.entries(factorMap)
    .map(([factor, count]) => count === 1 ? factor : `${factor}^${count}`)
    .join(' × ');
}

function fizzBuzz(start = 1, end = 100) {
  for (let i = start; i <= end; i++) {
    const factors = ` [${getPrimeFactors(i)}]`;
    
    if (i % 3 === 0 && i % 5 === 0) {
      console.log(`FizzBuzz${factors}`);
    } else if (i % 3 === 0) {
      console.log(`Fizz${factors}`);
    } else if (i % 5 === 0) {
      console.log(`Buzz${factors}`);
    } else {
      console.log(`${i}${factors}`);
    }
  }
}

// Execute the FizzBuzz function
console.log('Running FizzBuzz from 1 to 100 (ascending):');
fizzBuzz();