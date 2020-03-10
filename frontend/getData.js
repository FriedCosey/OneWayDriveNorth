$(function(){

    $.get("http://localhost:8080/sensors/microwave?starttime=1574614445408&endtime=1575151263833&status=1", function(data, status){
        console.log(data);
    });

    $.get("http://localhost:8080/sensors/microwave/doorCount?starttime=1574614445408&endtime=1575151263833&status=1", function(data, status){
        console.log(data);
    });

});
