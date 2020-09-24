module ALU (input signed [15:0] x,
            input signed [15:0] y,
            input zx,
            input nx,
            input zy,
            input ny,
            input f,
            input no,
            output reg signed [15:0] out,
            output reg zr,
            output reg ng);

reg signed [15:0] x_internal;
reg signed [15:0] y_internal;
reg signed [15:0] out_internal;

// x_internal statements computation
always @* begin
    if (zx)
        x_internal = 0;
    else
        x_internal = x;
        if (nx)
            x_internal = ~x_internal;
        else
            x_internal = x_internal;
end

// y_internal statements computation
always @* begin
    if (zy)
        y_internal = 0;
    else
        y_internal = y;
        if (ny)
            y_internal = ~y_internal;
        else
            y_internal = y_internal;
end

// out_internal statements computation
always @* begin
    if (f)
        out_internal = x_internal + y_internal;
    else
        out_internal = x_internal & y_internal;
end

// out statements computation
always @* begin
    if (no)
        out = ~out_internal;
    else
        out = out_internal;
end

// zr computation
always @* begin
    if (out == 0)
        zr = 1;
    else
        zr = 0;
end

// ng computation
always @* begin
    if (out < 0)
        ng = 1;
    else
        ng = 0;
end
endmodule // ALU
