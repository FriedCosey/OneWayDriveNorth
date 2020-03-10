function sensorData(id, width, height, data)
{
    var calData = randomData(width, height, data);
    console.log(calData);
    var grid = d3.select(id).append("svg")
                    .attr("width", width)
                    .attr("height", height)
                    .attr("class", "chart")
                    .attr("border-style", "dotted");

    var row = grid.selectAll(".row")
                  .data(calData)
                .enter().append("svg:g")
                  .attr("class", "row");

    var col = row.selectAll(".cell")
                 .data(function (d) { return d; })
                .enter().append("svg:rect")
                 .attr("class", "cell")
                 .attr("x", function(d) { return d.x; })
                 .attr("y", function(d) { return d.y; })
                 .attr("width", function(d) { return d.width; })
                 .attr("height", function(d) { return d.height; })
                 .attr("id", function(d) {return d.id})
                 .style("fill", '#FFF')
                 .style("stroke", '#555')

    var text = row.selectAll(".label")
        .data(function(d) {return d;})
      .enter().append("svg:text")
        .attr("x", function(d) { return d.x + d.width/2 })
        .attr("y", function(d) { return d.y + d.height/2 })
        .attr("text-anchor","middle")
        .attr("dy",".35em")
        .text(function(d) { return d.value });
}

function randomData(gridWidth, gridHeight, data)
{
    var res = new Array();
    var gridItemWidth = gridWidth / 8;
    var gridItemHeight =  gridHeight / 2;
    var startY = 0;
    var startX = 0;
    var stepX = gridItemWidth;
    var stepY = gridItemHeight;
    var xpos = startX;
    var ypos = startY;
    var days = ["ACTIVITY", "SUNDAY", "MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY"];
    var times = [2, 4, 5, 6, 7, 2, 3];
    var newValue;

    for (var index_a = 0; index_a < 2; index_a++)
    {
        res.push(new Array());
        for (var index_b = 0; index_b < 8; index_b++)
        {
            if (index_a == 0)
                newValue = days[index_b];
            else {
                newValue = 0;
                if (index_b == 0) {
            		newValue = "Open Microwave Door";
                }
                else {
                    newValue = data[index_b-1]["count"]
                }
            }
            res[index_a].push({
                                time: index_b,
                                value: newValue,
                                width: gridItemWidth,
                                height: gridItemHeight,
                                x: xpos,
                                y: ypos,
                                id: index_a * 8 + index_b
                            });
            xpos += stepX;
        }
        xpos = startX;
        ypos += stepY;
    }
    return res;
}


/*$.get("http://localhost:8080/sensors/microwave?starttime=1574614445408&endtime=1575151263833&status=1", function(data, status){
        console.log(data)
});*/
$.get("http://localhost:8080/sensors/microwave/doorCount?starttime=1574614445408&endtime=1575151263833&status=1", function(data, status){
    console.log(data);
    sensorData('#grid', 800, 100, data);
});
