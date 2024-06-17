# ffmpeg를 이용한 동영상 편집 서버
동영상 업로드, 동영상 Trim, 동영상 Concat, 동영상 정보 조회

# Reuslt Code
**0** : 요청 성공  
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

**postman example**  
<img width="522" alt="image" src="https://github.com/SundaePorkCutlet/video-edit/assets/87690981/67dc081d-8db1-44d4-9f08-89372d8ee177">


# API Example
### `:{port}/upload`  
동영상 업로드  
Multipart form 형식으로 업로드  

**postman example**
<img width="834" alt="image" src="https://github.com/SundaePorkCutlet/video-edit/assets/87690981/9ae267ef-1423-4a18-b6d2-a8f09df11319">
