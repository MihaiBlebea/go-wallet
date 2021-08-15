// https://gionkunz.github.io/chartist-js/getting-started.html

(async function() {
    const getData = async (start, end) => {
        try {
            let response = await axios.get('/api/v1/report?start=' + start + '&end=' + end)
            console.log(response)

            if (response.status !== 200) {
                throw Error('Request error, status not 200')
            }

            return response.data
        } catch(err) {
            console.error(err)

            return null
        }
    }

    let params = new URLSearchParams(window.location.search)

    let report = await getData(params.get('start'), params.get('end'))
    
    var data = {
        // A labels array that can contain any sort of values
        labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri'],
        // Our series array that contains series objects or in this case series data arrays
        series: [
            [5, 2, 4, 2, 0],
            [15, 2, 4, 20, 0]
        ]
    };
    
    // Create a new line chart object where as first parameter we pass in a selector
    // that is resolving to our chart container element. The Second parameter
    // is the actual data object.
    new Chartist.Line('.ct-chart', data);
})()