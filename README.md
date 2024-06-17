# ffmpeg를 이용한 동영상 편집 서버
동영상 업로드, 동영상 Trim, 동영상 Concat, 동영상 정보 조회

# Reuslt Code
**200** : 요청 성공  
**400** : 요청 오류  
**500** : 서버 오류  

**1001** : Trim 성공  
**1002** : Concat 성공
**1003** : Trim 성공 Concat 실패  
**1004** : Concat 성공 Trim 실패  
**1005** : Trim Concat 성공  

**6001** : 허용하는 확장자가 아님  
**6002** : Trim 실패  
**6003** : Concat 실패    

**example**  
  
<pre><code>{
    "code": 6002
}</code></pre>


# API Example
### `:{port}/upload`  
동영상 업로드  
Multipart form 형식으로 업로드  

**postman example**
  
<img width="834" alt="image" src="https://github.com/SundaePorkCutlet/video-edit/assets/87690981/9ae267ef-1423-4a18-b6d2-a8f09df11319">

### `:{port}/modify`
3가지 조건에 따라 parameter 변화  

### 1. trim만 요청하는 경우
<pre><code>{
    "isTrimed":true,
    "trimVideoList":[
        {
        "videoId":"2d05c690-ee19-49d2-b256-85a1641cd24b",
        "startTime":2,
        "endTime":5
        },
        {
        "videoId":"210a0287-2d83-4b80-afa5-44cdd709650d",
        "startTime":1,
        "endTime":70
        }
    ]
}</code></pre>

isTrimed : true로 요청
trimVideoList : 배열 형식으로 여러 동영상 trim 요청 가능  
videoId : 요청할 동영상 고유 id  
startTime : 시작 시간 (초 단위)  
endTime : 종료 시간 (초 단위)  

### 2. concat만 요청하는 경우
<pre><code>{
    "isConcated":true,
    "concatVideoList":[
        "0c649117-6ae9-4d22-b033-80c646191fa5","c801f057-d346-4889-8d5d-68878423f03e"
    ]
}</code></pre>  

  isConcated : true로 요청  
  concatVideoList : 요청할 동영상 고유 id 리스트  

### 3. trim한 동영상들을 concat 요청하는 경우  
<pre><code>{
    "isTrimed":true,
    "trimVideoList":[
         {
        "videoId":"7d04454e-5fd9-48b8-a8bd-ebdabd7e3362",
        "startTime":2,
        "endTime":5
        },
        {
        "videoId":"83158ed9-5ff4-46b2-aed9-63064196b505",
        "startTime":1,
        "endTime":7
        }
    ],
    "isConcated":true
}</code></pre>  

trim과 똑같은 요청에서 isConcated만 true요청 추가
