# gobench

Descripción de Cada Dato:
name:

Tipo: string
Descripción: Es el nombre de la prueba. En este caso, el valor es "Local", lo que podría indicar que la prueba se realizó en un entorno local o en un servidor local.
method:

Tipo: string
Descripción: Es el método HTTP que se utilizó en la prueba. En este caso, el valor es "GET", lo que indica que se realizaron solicitudes GET al servidor.
requests:

Tipo: int
Descripción: Es el número total de solicitudes realizadas durante la prueba. En este caso, se hicieron 10 solicitudes.
duration:

Tipo: float64
Descripción: Representa la duración total de la prueba en segundos. En este caso, la duración total fue 0.0321 segundos. Este valor mide cuánto tiempo en total se tomó para completar las 10 solicitudes.
rps (Requests per Second):

Tipo: float64
Descripción: Mide el número de solicitudes procesadas por segundo durante la prueba. En este caso, 311.45 solicitudes por segundo. Este valor se calcula dividiendo el número total de solicitudes entre la duración de la prueba (en segundos).
avg_latency:

Tipo: int
Descripción: Es la latencia promedio de las solicitudes en microsegundos (usualmente se mide en microsegundos, pero depende de la implementación). En este caso, la latencia promedio es 20,732,222 microsegundos (aproximadamente 20.7 segundos), lo que es un valor bastante alto, por lo que podría ser un error o un valor inusualmente alto debido a algún problema en la prueba o en el servidor.
max_latency:

Tipo: int
Descripción: Representa la latencia máxima experimentada en alguna de las solicitudes, medida en microsegundos. En este caso, el valor es 31,989,805 microsegundos (~32 segundos), lo que es muy alto y podría indicar un gran retraso o un cuellos de botella en el servidor al que se realizaron las solicitudes.
min_latency:

Tipo: int
Descripción: Es la latencia mínima de las solicitudes, medida en microsegundos. En este caso, la latencia mínima fue 6,048,423 microsegundos (~6 segundos). Esto muestra que algunas solicitudes fueron mucho más rápidas, pero aún así, es un valor relativamente alto.


Resumen:
name: El nombre de la prueba (en este caso "Local").
method: El método HTTP utilizado (en este caso GET).
requests: Número total de solicitudes realizadas en la prueba.
duration: Duración total de la prueba en segundos.
rps: Solicitudes por segundo (requests per second).
avg_latency: Latencia promedio de las solicitudes en microsegundos.
max_latency: Latencia máxima de las solicitudes en microsegundos.
min_latency: Latencia mínima de las solicitudes en microsegundos.
