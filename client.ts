import net from 'net'

let conn = net.connect({
  host: 'localhost',
  port: 8500,
})

function set(key: string, value: string | number) {
  conn.write(`set ${key} ${value}\n`)
}

function rename(oldKey: string, newKey: string | number) {
  conn.write(`rename ${oldKey} ${newKey}\n`)
}

function del(key: string) {
  conn.write(`del ${key}\n`)
}

function get(key: string) {
  return new Promise<string>((resolve, reject) => {
    conn.write(`get ${key}\n`)
    conn.once('data', chunk => {
      let text = chunk.toString()
      if (text === '(none)\n') {
        reject('key not found')
      } else {
        let value = text.slice('value:'.length, text.length - 1)
        resolve(value)
      }
    })
  })
}

async function main() {
  try {
    let oldValue = await get('alice')
    let newValue = +oldValue + 1
    set('alice', newValue)
  } catch (error) {
    set('alice', 10)
  }

  let value = await get('alice')
  console.log('value:', value)
}
main().catch(e => console.error(e))

// conn.write('set alice 10\n')
// conn.write('get alice\n')

// conn.on('data', chunk => {
//   console.log('received:', chunk.toString('ascii'))
// })
