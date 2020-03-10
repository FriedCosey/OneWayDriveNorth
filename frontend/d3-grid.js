function sensorData(id, width, height, square)
{
    var calData = randomData(width, height, square);
    console.log(calData);
    var grid = d3.select(id).append("svg")
                    .attr("width", width)
                    .attr("height", height)
                    .attr("class", "chart");

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
                 .on('mouseover', function() {
                    d3.select(this)
                        .style('fill', '#0F0');
                 })
                 .on('mouseout', function() {
                    d3.select(this)
                        .style('fill', '#FFF');
                 })
                 .on('click', function() {
                    console.log(d3.select(this));
                 })
                 .style("fill", '#FFF')
                 .style("stroke", '#555')

   /* var id = 0;

    for (var i = 1; i < 2; i++) {
        for(var q = 0; q < 8; q++) {
            id =
        }
    }
    */
    var text = row.selectAll(".label")
        .data(function(d) {return d;})
      .enter().append("svg:text")
        .attr("x", function(d) { return d.x + d.width/2 })
        .attr("y", function(d) { return d.y + d.height/2 })
        .attr("text-anchor","middle")
        .attr("dy",".35em")
        .text(function(d) { return d.value });
}

function randomData(gridWidth, gridHeight)
{
    var data = new Array();
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
        data.push(new Array());
        for (var index_b = 0; index_b < 8; index_b++)
        {
            if (index_a == 0)
                newValue = days[index_b];
            else {
                newValue = 0;
                if (index_b == 0) {
            		newValue = "Open Microwave Door";
            	}
            }
            data[index_a].push({
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
        //console.log(data[0])
        xpos = startX;
        ypos += stepY;
    }
    return data;
}
sensorData('#grid', 1400, 100);
$.get("http://localhost:8080/sensors/microwave?starttime=1574614445408&endtime=1575151263833&status=1", function(data, status){
        console.log(data);
});
